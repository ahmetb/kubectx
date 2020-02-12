#!/usr/bin/perl
#

my %cluster;
my %server;
my %clientcrt;
my %clientkey;

my $outfile;

if ( !-e $ENV{"HOME"} . "/.kube/contexts" ) {
    print "yes ,it exists!" ;
    die;
}

opendir(DH, $ENV{"HOME"} . "/.kube/contexts") or die $!;
my @files = readdir(DH);
closedir(DH);

#looping trough files
foreach my $fi (@files){ 
  if ($fi =~ /(.*)\.conf/){
    print "$1 - $fi\n";
    my $filen = $1;
    open(FH, '<', $ENV{"HOME"} . "/.kube/contexts/$fi") or die $!;
    my @lines = <FH>;
    close(FH);

    foreach my $value (@lines){
       if ($value =~ /certificate-authority-data:/){
	       # 	 print $value;
         $cluster{$filen} = $value;
       }
       if ($value =~ /server:/){
    	   #     print $value;
         $server{$filen} = $value;
       }
       if ($value =~ /client-certificate-data:/){
    	   #     print $value;
         $clientcrt{$filen} = $value;
       }
       if ($value =~ /client-key-data:/){
    	   #     print $value;
         $clientkey{$filen} = $value;
       }
    }
    
  }

}



$outfile = "apiVersion: v1\nclusters:\n";

foreach my $fi (@files){
  if ($fi =~ /(.*)\.conf/){
    my $filen = $1;
    $outfile = $outfile . "- cluster:\n";
    $outfile = $outfile . $cluster{$filen};
    $outfile = $outfile . $server{$filen};
    $outfile = $outfile . "  name: " . $filen ."\n";
  }
}

$outfile = $outfile . "contexts:\n";

foreach my $fi (@files){
  if ($fi =~ /(.*)\.conf/){
    my $filen = $1; 
    $outfile = $outfile . "- context:\n";
    $outfile = $outfile . "    cluster: " . $filen ."\n";
    $outfile = $outfile . "    user: " . $filen . "-usr\n";
    $outfile = $outfile . "  name: " . $filen ."\n";

  }
}

$outfile = $outfile . "kind: Config\n";
$outfile = $outfile . "preferences: {}\n";
$outfile = $outfile . "users:\n";

foreach my $fi (@files){
  if ($fi =~ /(.*)\.conf/){
    my $filen = $1; 
    $outfile = $outfile . "- name: " . $filen . "-usr\n";
    $outfile = $outfile . "  user: \n";
    $outfile = $outfile . $clientcrt{$filen};
    $outfile = $outfile . $clientkey{$filen};
  }
}




#print $outfile;
my $filename = $ENV{"HOME"} . "/.kube/config";
open(my $fh, '>', $filename) or die "Could not open file '$filename' $!";
print $fh $outfile;
close $fh;
